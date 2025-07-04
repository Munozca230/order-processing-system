package com.example.orderworker.service;

import com.example.orderworker.model.CustomerDetails;
import com.example.orderworker.model.OrderMessage;
import com.example.orderworker.model.ProductDetails;
import com.example.orderworker.repository.FailedMessageRepository;
import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.test.StepVerifier;

import java.util.List;

import static com.github.tomakehurst.wiremock.client.WireMock.*;
import static org.mockito.Mockito.mock;

class EnrichmentServiceTest {

    private static WireMockServer wireMockServer;
    private static EnrichmentService enrichmentService;

    @BeforeAll
    static void setup() {
        wireMockServer = new WireMockServer(WireMockConfiguration.options().dynamicPort());
        wireMockServer.start();
        configureFor("localhost", wireMockServer.port());

        WebClient testClient = WebClient.builder().baseUrl("http://localhost:" + wireMockServer.port()).build();
        
        // Mock the RetryService for testing
        FailedMessageRepository mockFailedMessageRepository = mock(FailedMessageRepository.class);
        RetryService mockRetryService = new RetryService(mockFailedMessageRepository);
        
        // Use same client for both product and customer for simplicity
        enrichmentService = new EnrichmentService(testClient, testClient, mockRetryService);
    }

    @AfterAll
    static void tearDown() {
        wireMockServer.stop();
    }

    @Test
    @DisplayName("Debe enriquecer la orden con customer y products")
    void enrich_shouldReturnEnrichedOrder() {
        // Stubs
        stubFor(get(urlEqualTo("/customers/c1"))
                .willReturn(okJson("{\"customerId\":\"c1\",\"name\":\"John\",\"active\":true}")));

        stubFor(get(urlEqualTo("/products/p1"))
                .willReturn(okJson("{\"productId\":\"p1\",\"name\":\"Widget\",\"price\":9.99}")));

        OrderMessage order = new OrderMessage("o1", "c1", List.of("p1"));

        StepVerifier.create(enrichmentService.enrich(order))
                .expectNextMatches(enriched ->
                        enriched.customer().customerId().equals("c1") &&
                        enriched.products().size() == 1 &&
                        enriched.products().get(0).productId().equals("p1"))
                .verifyComplete();
    }
}
