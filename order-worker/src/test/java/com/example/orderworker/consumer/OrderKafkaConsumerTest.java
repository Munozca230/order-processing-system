package com.example.orderworker.consumer;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class OrderKafkaConsumerTest {

    @Test
    @DisplayName("El consumidor debería procesar el mensaje sin lanzar excepciones")
    void consume_shouldNotThrow() {
        OrderKafkaConsumer consumer = new OrderKafkaConsumer();
        assertDoesNotThrow(() -> consumer.consume("{\"orderId\":\"test-1\"}"));
    }
}
