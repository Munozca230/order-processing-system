package com.example.orderworker.service;

import com.example.orderworker.exception.InvalidOrderException;
import com.example.orderworker.model.EnrichedOrder;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

@Service
public class ValidationService {

    public Mono<EnrichedOrder> validate(EnrichedOrder order) {
        if (order.customer() == null || !order.customer().active()) {
            return Mono.error(new InvalidOrderException("Customer is inactive or not found"));
        }
        if (order.products() == null || order.products().isEmpty()) {
            return Mono.error(new InvalidOrderException("No products found for order"));
        }
        // extra rules can be added here
        return Mono.just(order);
    }
}
