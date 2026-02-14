package com.efact.validator.util;

import com.efact.validator.model.Documento;
import com.fasterxml.jackson.databind.ObjectMapper;

public final class DocumentUtils {

    private DocumentUtils() {
        throw new UnsupportedOperationException();
    }

    public static Documento cloneDocument(Documento original, ObjectMapper objectMapper) {
        try {
            String json = objectMapper.writeValueAsString(original);
            return objectMapper.readValue(json, Documento.class);
        } catch (Exception e) {
            throw new RuntimeException("Error al clonar documento", e);
        }
    }
}
