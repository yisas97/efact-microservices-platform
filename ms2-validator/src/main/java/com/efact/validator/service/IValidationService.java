package com.efact.validator.service;

import com.efact.validator.model.Documento;

public interface IValidationService {

    boolean validateDocument(Documento document);
}
