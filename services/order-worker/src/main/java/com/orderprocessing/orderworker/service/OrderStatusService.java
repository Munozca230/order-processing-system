package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.dto.OrderStatusResponse;
import com.orderprocessing.orderworker.model.FailedMessage;
import com.orderprocessing.orderworker.model.OrderDocument;
import com.orderprocessing.orderworker.model.OrderStatus;
import com.orderprocessing.orderworker.repository.FailedMessageRepository;
import com.orderprocessing.orderworker.repository.OrderRepository;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.time.LocalDateTime;
import java.time.ZoneOffset;

@Service
public class OrderStatusService {
    
    private final OrderRepository orderRepository;
    private final FailedMessageRepository failedMessageRepository;
    
    public OrderStatusService(OrderRepository orderRepository, FailedMessageRepository failedMessageRepository) {
        this.orderRepository = orderRepository;
        this.failedMessageRepository = failedMessageRepository;
    }
    
    public Mono<OrderStatusResponse> getOrderStatus(String orderId) {
        return orderRepository.findByOrderId(orderId)
            .map(this::mapToCompletedStatus)
            .switchIfEmpty(checkFailedStatus(orderId))
            .switchIfEmpty(Mono.just(OrderStatusResponse.processing(orderId)));
    }
    
    private Mono<OrderStatusResponse> checkFailedStatus(String orderId) {
        // Try different key patterns that could exist in Redis
        return failedMessageRepository.getFailedMessage(orderId)
            .switchIfEmpty(failedMessageRepository.getFailedMessage("order_" + orderId))
            .map(this::mapToFailedStatus);
    }
    
    private OrderStatusResponse mapToCompletedStatus(OrderDocument order) {
        return new OrderStatusResponse(
            order.orderId(),
            order.status().name(),
            order.createdAt(),
            order.processedAt(),
            order.updatedAt(),
            order.failureReason(),
            order.retryCount(),
            null
        );
    }
    
    private OrderStatusResponse mapToFailedStatus(FailedMessage failedMessage) {
        OrderStatus status = failedMessage.getRetryCount() > 0 && 
                           failedMessage.getNextRetryAt().isAfter(LocalDateTime.now()) 
                           ? OrderStatus.RETRYING : OrderStatus.FAILED;
                           
        return new OrderStatusResponse(
            failedMessage.getMessageId(),
            status.name(),
            failedMessage.getLastRetryAt(),
            null,
            failedMessage.getLastRetryAt(),
            failedMessage.getErrorMessage(),
            failedMessage.getRetryCount(),
            failedMessage.getNextRetryAt()
        );
    }
}