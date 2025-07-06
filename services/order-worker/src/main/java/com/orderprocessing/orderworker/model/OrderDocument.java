package com.orderprocessing.orderworker.model;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDateTime;
import java.util.List;

@Document(collection = "orders")
public record OrderDocument(
        @Id String id,
        String orderId,
        String customerId,
        List<ProductDetails> products,
        OrderStatus status,
        LocalDateTime createdAt,
        LocalDateTime processedAt,
        LocalDateTime updatedAt,
        String failureReason,
        int retryCount
) {
    public static OrderDocument from(EnrichedOrder enriched) {
        return new OrderDocument(
                null,
                enriched.order().orderId(),
                enriched.order().customerId(),
                enriched.products(),
                OrderStatus.COMPLETED,
                LocalDateTime.now(),
                LocalDateTime.now(),
                LocalDateTime.now(),
                null,
                0
        );
    }
    
    public static OrderDocument processing(String orderId, String customerId, List<ProductDetails> products) {
        return new OrderDocument(
                null,
                orderId,
                customerId,
                products,
                OrderStatus.PROCESSING,
                LocalDateTime.now(),
                null,
                LocalDateTime.now(),
                null,
                0
        );
    }
}
