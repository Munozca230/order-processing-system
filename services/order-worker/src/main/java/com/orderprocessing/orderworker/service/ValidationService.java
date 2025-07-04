package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.exception.InvalidOrderException;
import com.orderprocessing.orderworker.model.EnrichedOrder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

@Service
public class ValidationService {

    private static final Logger logger = LoggerFactory.getLogger(ValidationService.class);

    public Mono<EnrichedOrder> validate(EnrichedOrder order) {
        logger.info("🔍 VALIDATION START for order: {}", order.order().orderId());
        
        if (order.customer() == null) {
            logger.error("❌ VALIDATION FAILED: Customer is null for order: {}", order.order().orderId());
            return Mono.error(new InvalidOrderException("Customer is null"));
        }
        
        if (!order.customer().active()) {
            logger.error("❌ VALIDATION FAILED: Customer {} is inactive for order: {}", 
                order.customer().customerId(), order.order().orderId());
            return Mono.error(new InvalidOrderException("Customer is inactive"));
        }
        
        if (order.products() == null || order.products().isEmpty()) {
            logger.error("❌ VALIDATION FAILED: No products found for order: {}", order.order().orderId());
            return Mono.error(new InvalidOrderException("No products found for order"));
        }
        
        logger.info("✅ VALIDATION SUCCESS for order: {} - customer: {}, products: {}", 
            order.order().orderId(), order.customer().customerId(), order.products().size());
        
        return Mono.just(order);
    }
}
