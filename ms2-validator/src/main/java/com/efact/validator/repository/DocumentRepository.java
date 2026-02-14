package com.efact.validator.repository;

import com.efact.validator.model.Documento;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface DocumentRepository extends MongoRepository<Documento, String> {
    Optional<Documento> findByIdDocumento(String idDocumento);
}
