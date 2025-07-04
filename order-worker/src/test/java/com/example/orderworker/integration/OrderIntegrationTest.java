package com.example.orderworker.integration;

import com.example.orderworker.event.ProcessedOrderEvent;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.ApplicationListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.KafkaContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;
import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import static com.github.tomakehurst.wiremock.client.WireMock.*;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Testcontainers
@SpringBootTest
class OrderIntegrationTest {

    @Container
    static KafkaContainer kafka = new KafkaContainer(DockerImageName.parse("confluentinc/cp-kafka:7.5.0"));

    private static WireMockServer productMock;
    private static WireMockServer customerMock;

    private static final CountDownLatch latch = new CountDownLatch(1);

    @DynamicPropertySource
    static void properties(DynamicPropertyRegistry registry) {
        productMock = new WireMockServer(WireMockConfiguration.options().dynamicPort());
        customerMock = new WireMockServer(WireMockConfiguration.options().dynamicPort());
        productMock.start();
        customerMock.start();

        registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
        registry.add("app.product-api.base-url", () -> "http://localhost:" + productMock.port());
        registry.add("app.customer-api.base-url", () -> "http://localhost:" + customerMock.port());
    }

    @BeforeAll
    static void stubs() {
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
    void registerListener(org.springframework.context.ApplicationEventPublisher publisher) {
        publisher.addApplicationListener((ApplicationListener<ProcessedOrderEvent>) event -> latch.countDown());
    }

    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;

    @Test
    @DisplayName("Debe procesar un mensaje de orden end-to-end")
    void shouldProcessOrder() throws InterruptedException {
        String msg = "{\"orderId\":\"o1\",\"customerId\":\"c1\",\"products\":[\"p1\"]}";
        kafkaTemplate.send(new ProducerRecord<>("orders", null, msg));

        boolean processed = latch.await(10, TimeUnit.SECONDS);
        assertTrue(processed, "La orden no fue procesada a tiempo");
    }
}
