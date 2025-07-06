package com.orderprocessing.orderworker.integration;

import com.orderprocessing.orderworker.event.ProcessedOrderEvent;
import com.orderprocessing.orderworker.repository.OrderRepository;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.junit.jupiter.api.AfterAll;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.ApplicationListener;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.KafkaContainer;
import org.testcontainers.containers.MongoDBContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;
import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;

import org.apache.kafka.clients.admin.AdminClient;
import org.apache.kafka.clients.admin.AdminClientConfig;
import org.apache.kafka.clients.admin.NewTopic;
import org.testcontainers.containers.wait.strategy.Wait;

import java.time.Duration;
import java.util.List;
import java.util.Map;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import static com.github.tomakehurst.wiremock.client.WireMock.*;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Testcontainers(disabledWithoutDocker = true)
@SpringBootTest
class OrderLockIntegrationTest {

    @Container
    static KafkaContainer kafka = new KafkaContainer(DockerImageName.parse("confluentinc/cp-kafka:7.5.0"))
            .withStartupTimeout(Duration.ofSeconds(60));

    @Container
    static GenericContainer<?> redis = new GenericContainer<>(DockerImageName.parse("redis:7-alpine"))
            .withExposedPorts(6379);

    @Container
    static MongoDBContainer mongo = new MongoDBContainer(DockerImageName.parse("mongo:7.0"))
            .waitingFor(Wait.forLogMessage(".*Waiting for connections.*\\n", 1))
            .withStartupTimeout(Duration.ofSeconds(60));

    private static WireMockServer productMock;
    private static WireMockServer customerMock;

    private static final CountDownLatch latch = new CountDownLatch(1);

    @DynamicPropertySource
    static void properties(DynamicPropertyRegistry registry) {
        productMock = new WireMockServer(WireMockConfiguration.options().dynamicPort());
        customerMock = new WireMockServer(WireMockConfiguration.options().dynamicPort());
        productMock.start();
        customerMock.start();

        // Ensure 'orders' topic exists before the application context starts
        try (AdminClient admin = AdminClient.create(
                Map.of(AdminClientConfig.BOOTSTRAP_SERVERS_CONFIG, kafka.getBootstrapServers()))) {
            admin.createTopics(List.of(new NewTopic("orders", 1, (short) 1))).all().get();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
        registry.add("spring.kafka.consumer.auto-offset-reset", () -> "earliest");
        registry.add("spring.data.mongodb.uri", mongo::getReplicaSetUrl);
        registry.add("spring.data.mongodb.database", () -> "test");
        registry.add("spring.redis.host", () -> redis.getHost());
        registry.add("spring.redis.port", () -> redis.getMappedPort(6379).toString());
        registry.add("app.product-api.base-url", () -> "http://localhost:" + productMock.port());
        registry.add("app.customer-api.base-url", () -> "http://localhost:" + customerMock.port());

        // Stubs
        productMock.stubFor(get(urlEqualTo("/products/p1"))
                .willReturn(okJson("{\"productId\":\"p1\",\"name\":\"Widget\",\"price\":9.99}")));
        customerMock.stubFor(get(urlEqualTo("/customers/c1"))
                .willReturn(okJson("{\"customerId\":\"c1\",\"name\":\"John\",\"active\":true}")));
    }

    @AfterAll
    static void shutdown() {
        productMock.stop();
        customerMock.stop();
    }

    @Autowired
    void registerListener(ConfigurableApplicationContext context) {
        context.addApplicationListener((ApplicationListener<ProcessedOrderEvent>) event -> latch.countDown());
    }

    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;

    @Autowired
    private OrderRepository orderRepository;

    @Test
    @DisplayName("No debe persistir duplicados si llegan dos mensajes con el mismo orderId")
    void shouldProcessOnlyOnce() throws InterruptedException {
        String msg = "{\"orderId\":\"dup-1\",\"customerId\":\"c1\",\"products\":[{\"productId\":\"p1\"}]}";

        // Enviar dos veces la misma orden
        kafkaTemplate.send(new ProducerRecord<>("orders", null, msg));
        kafkaTemplate.flush();
        kafkaTemplate.send(new ProducerRecord<>("orders", null, msg));
        kafkaTemplate.flush();

        boolean processed = latch.await(60, TimeUnit.SECONDS);
        assertTrue(processed, "La orden no fue procesada la primera vez");

        // Esperar un poco para asegurar que un posible segundo procesamiento termine
        TimeUnit.SECONDS.sleep(2);

        long count = orderRepository.count().block();
        assertEquals(1, count, "Se detectó más de una inserción de la misma orden");
    }
}
