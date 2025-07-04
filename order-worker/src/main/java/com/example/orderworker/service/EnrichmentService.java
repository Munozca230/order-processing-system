package com.example.orderworker.service;

import com.example.orderworker.model.*;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
public class EnrichmentService {

    private final WebClient productApiClient;
    private final WebClient customerApiClient;

    public EnrichmentService(@Qualifier("productApiClient") WebClient productApiClient,
                             @Qualifier("customerApiClient") WebClient customerApiClient) {
        this.productApiClient = productApiClient;
        this.customerApiClient = customerApiClient;
    }

    public Mono<EnrichedOrder> enrich(OrderMessage order) {
        Mono<CustomerDetails> customerMono = customerApiClient.get()
                .uri("/customers/{id}", order.customerId())
                .retrieve()
                .bodyToMono(CustomerDetails.class);

        Flux<ProductDetails> productsFlux = Flux.fromIterable(order.products())
                .flatMap(id -> productApiClient.get()
                        .uri("/products/{id}", id)
                        .retrieve()
                        .bodyToMono(ProductDetails.class));

        return Mono.zip(customerMono, productsFlux.collectList())
                .map(tuple -> new EnrichedOrder(order, tuple.getT1(), tuple.getT2()));
    }
}
