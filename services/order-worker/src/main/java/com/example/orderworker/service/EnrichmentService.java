package com.example.orderworker.service;

import com.example.orderworker.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
public class EnrichmentService {

    private static final Logger logger = LoggerFactory.getLogger(EnrichmentService.class);

    private final WebClient productApiClient;
    private final WebClient customerApiClient;
    private final RetryService retryService;

    public EnrichmentService(@Qualifier("productApiClient") WebClient productApiClient,
                             @Qualifier("customerApiClient") WebClient customerApiClient,
                             RetryService retryService) {
        this.productApiClient = productApiClient;
        this.customerApiClient = customerApiClient;
        this.retryService = retryService;
    }

    public Mono<EnrichedOrder> enrich(OrderMessage order) {
        String messageId = "order_" + order.orderId();
        String messageContent = order.toString();

        return retryService.executeWithRetry(messageId, messageContent, () -> {
            logger.info("Enriching order: {}", order.orderId());

            Mono<CustomerDetails> customerMono = fetchCustomerWithRetry(order.customerId(), messageId);
            Flux<ProductDetails> productsFlux = fetchProductsWithRetry(order.products(), messageId);

            return Mono.zip(customerMono, productsFlux.collectList())
                    .map(tuple -> new EnrichedOrder(order, tuple.getT1(), tuple.getT2()))
                    .doOnSuccess(enrichedOrder -> logger.info("Successfully enriched order: {}", order.orderId()))
                    .doOnError(error -> logger.error("Failed to enrich order: {}", order.orderId(), error));
        });
    }

    private Mono<CustomerDetails> fetchCustomerWithRetry(String customerId, String messageId) {
        return retryService.executeWithRetry(
            messageId + "_customer_" + customerId,
            "customer:" + customerId,
            () -> customerApiClient.get()
                .uri("/customers/{id}", customerId)
                .retrieve()
                .bodyToMono(CustomerDetails.class)
                .doOnNext(customer -> logger.debug("Fetched customer: {}", customerId))
        );
    }

    private Flux<ProductDetails> fetchProductsWithRetry(java.util.List<String> productIds, String messageId) {
        return Flux.fromIterable(productIds)
                .flatMap(productId -> retryService.executeWithRetry(
                    messageId + "_product_" + productId,
                    "product:" + productId,
                    () -> productApiClient.get()
                        .uri("/products/{id}", productId)
                        .retrieve()
                        .bodyToMono(ProductDetails.class)
                        .doOnNext(product -> logger.debug("Fetched product: {}", productId))
                ));
    }
}
