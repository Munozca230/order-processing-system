package com.example.orderworker.model;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.Instant;
import java.util.List;

@Document(collection = "orders")
public record OrderDocument(
        @Id String id,
        String orderId,
        String customerId,
        List<ProductDetails> products
) {
    public static OrderDocument from(EnrichedOrder enriched) {
        return new OrderDocument(
                null,
                enriched.order().orderId(),
                enriched.order().customerId(),
                enriched.products()
        );
    }
}
