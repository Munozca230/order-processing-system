package com.orderprocessing.orderworker.repository;

import com.orderprocessing.orderworker.model.OrderDocument;
import org.springframework.data.mongodb.repository.ReactiveMongoRepository;
import reactor.core.publisher.Mono;

public interface OrderRepository extends ReactiveMongoRepository<OrderDocument, String> {
    Mono<OrderDocument> findByOrderId(String orderId);
}
