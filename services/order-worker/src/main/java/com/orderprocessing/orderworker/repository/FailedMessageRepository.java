package com.orderprocessing.orderworker.repository;

import com.orderprocessing.orderworker.model.FailedMessage;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import org.springframework.data.domain.Range;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.time.LocalDateTime;

@Repository
public class FailedMessageRepository {

    private final ReactiveRedisTemplate<String, String> redisTemplate;
    private final ObjectMapper objectMapper;
    private static final String FAILED_MESSAGES_KEY = "failed_messages:";
    private static final String RETRY_QUEUE_KEY = "retry_queue";

    public FailedMessageRepository(ReactiveRedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
        this.objectMapper = new ObjectMapper();
        this.objectMapper.registerModule(new JavaTimeModule());
    }

    public Mono<Void> saveFailedMessage(FailedMessage failedMessage) {
        return Mono.fromCallable(() -> {
            try {
                return objectMapper.writeValueAsString(failedMessage);
            } catch (JsonProcessingException e) {
                throw new RuntimeException("Error serializing failed message", e);
            }
        })
        .flatMap(json -> redisTemplate.opsForValue()
            .set(FAILED_MESSAGES_KEY + failedMessage.getMessageId(), json, Duration.ofDays(7)))
        .then();
    }

    public Mono<FailedMessage> getFailedMessage(String messageId) {
        return redisTemplate.opsForValue()
            .get(FAILED_MESSAGES_KEY + messageId)
            .map(json -> {
                try {
                    return objectMapper.readValue(json, FailedMessage.class);
                } catch (JsonProcessingException e) {
                    throw new RuntimeException("Error deserializing failed message", e);
                }
            });
    }

    public Mono<Void> deleteFailedMessage(String messageId) {
        return redisTemplate.opsForValue()
            .delete(FAILED_MESSAGES_KEY + messageId)
            .then();
    }

    public Mono<Void> addToRetryQueue(String messageId, LocalDateTime nextRetryTime) {
        long score = nextRetryTime.atZone(java.time.ZoneOffset.UTC).toEpochSecond();
        return redisTemplate.opsForZSet()
            .add(RETRY_QUEUE_KEY, messageId, score)
            .then();
    }

    public Flux<String> getMessagesReadyForRetry() {
        long currentTime = LocalDateTime.now().atZone(java.time.ZoneOffset.UTC).toEpochSecond();
        return redisTemplate.opsForZSet()
            .rangeByScore(RETRY_QUEUE_KEY, Range.closed(0.0, (double) currentTime));
    }

    public Mono<Void> removeFromRetryQueue(String messageId) {
        return redisTemplate.opsForZSet()
            .remove(RETRY_QUEUE_KEY, messageId)
            .then();
    }
}