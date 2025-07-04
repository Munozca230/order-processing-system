package com.orderprocessing.orderworker.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.ReactiveStringRedisTemplate;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.time.Duration;

/**
 * Servicio que implementa un lock simple por orderId usando Redis.
 * Usa SETNX con TTL para evitar reprocesos mientras la orden se procesa.
 */
@Service
public class OrderLockService {

    private final ReactiveStringRedisTemplate redisTemplate;
    private final Duration ttl;

    public OrderLockService(ReactiveStringRedisTemplate redisTemplate,
                            @Value("${app.order-lock.ttl-seconds:30}") long ttlSeconds) {
        this.redisTemplate = redisTemplate;
        this.ttl = Duration.ofSeconds(ttlSeconds);
    }

    /**
     * Intenta adquirir el lock. Devuelve true si lo obtuvo, false si ya estaba tomado.
     */
    public Mono<Boolean> acquire(String orderId) {
        String key = "order:lock:" + orderId;
        return redisTemplate
                .opsForValue()
                .setIfAbsent(key, "1", ttl)
                .defaultIfEmpty(false);
    }

    /**
     * Libera el lock antes de su expiración.
     */
    public Mono<Long> release(String orderId) {
        String key = "order:lock:" + orderId;
        return redisTemplate.delete(key);
    }
}
