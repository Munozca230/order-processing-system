package com.orderprocessing.orderworker.repository;

import com.orderprocessing.orderworker.model.OrderDocument;
import org.springframework.data.mongodb.repository.ReactiveMongoRepository;

public interface OrderRepository extends ReactiveMongoRepository<OrderDocument, String> {
    // Additional query methods if needed
}
