package com.example.orderworker.repository;

import com.example.orderworker.model.OrderDocument;
import org.springframework.data.mongodb.repository.ReactiveMongoRepository;

public interface OrderRepository extends ReactiveMongoRepository<OrderDocument, String> {
    // Additional query methods if needed
}
