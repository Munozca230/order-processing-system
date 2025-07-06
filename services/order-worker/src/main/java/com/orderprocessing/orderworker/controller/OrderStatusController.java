package com.orderprocessing.orderworker.controller;

import com.orderprocessing.orderworker.dto.OrderStatusResponse;
import com.orderprocessing.orderworker.service.OrderStatusService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/orders")
public class OrderStatusController {
    
    private final OrderStatusService orderStatusService;
    
    public OrderStatusController(OrderStatusService orderStatusService) {
        this.orderStatusService = orderStatusService;
    }
    
    @GetMapping("/{orderId}/status")
    public Mono<ResponseEntity<OrderStatusResponse>> getOrderStatus(@PathVariable String orderId) {
        return orderStatusService.getOrderStatus(orderId)
            .map(ResponseEntity::ok)
            .defaultIfEmpty(ResponseEntity.notFound().build());
    }
}