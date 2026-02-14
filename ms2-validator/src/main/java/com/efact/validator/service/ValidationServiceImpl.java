package com.efact.validator.service;

import com.efact.validator.constants.MessageConstants;
import com.efact.validator.model.Documento;
import com.efact.validator.model.Item;
import com.efact.validator.util.MathUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import static com.efact.validator.constants.ValidationConstants.IGV_RATE;
import static com.efact.validator.constants.ValidationConstants.TOLERANCE;

@Service
public class ValidationServiceImpl implements IValidationService {

    private static final Logger logger = LoggerFactory.getLogger(ValidationServiceImpl.class);

    @Override
    public boolean validateDocument(Documento document) {
        logger.info("Validando documento: {}", document.getIdDocumento());

        if (document.getItems() == null || document.getItems().isEmpty()) {
            logger.error(MessageConstants.DOCUMENTO_SIN_ITEMS);
            return false;
        }

        if (!validateItems(document)) {
            logger.error(MessageConstants.VALIDACION_ITEMS_FALLIDA);
            return false;
        }

        if (!validateTotals(document)) {
            logger.error(MessageConstants.VALIDACION_TOTALES_FALLIDA);
            return false;
        }

        logger.info("Validación del documento exitosa");
        return true;
    }

    private boolean validateItems(Documento document) {
        for (int i = 0; i < document.getItems().size(); i++) {
            Item item = document.getItems().get(i);

            double expectedPrecioTotal = item.getPrecioUnitario() * item.getCantidad();
            double expectedIgvTotal = expectedPrecioTotal * IGV_RATE;

            if (!MathUtils.areEqual(item.getPrecioTotal(), expectedPrecioTotal, TOLERANCE)) {
                logger.error("Ítem {}: precioTotal no coincide. Esperado: {}, Obtenido: {}",
                    i, expectedPrecioTotal, item.getPrecioTotal());
                return false;
            }

            if (!MathUtils.areEqual(item.getIgvTotal(), expectedIgvTotal, TOLERANCE)) {
                logger.error("Ítem {}: igvTotal no coincide. Esperado: {}, Obtenido: {}",
                    i, expectedIgvTotal, item.getIgvTotal());
                return false;
            }
        }
        return true;
    }

    private boolean validateTotals(Documento document) {
        double expectedMontoSinImpuestos = document.getItems().stream()
            .mapToDouble(Item::getPrecioTotal)
            .sum();

        double expectedIgvTotal = expectedMontoSinImpuestos * IGV_RATE;
        double expectedMontoTotal = expectedMontoSinImpuestos + expectedIgvTotal;

        if (!MathUtils.areEqual(document.getMontoTotalSinImpuestos(), expectedMontoSinImpuestos, TOLERANCE)) {
            logger.error("montoTotalSinImpuestos no coincide. Esperado: {}, Obtenido: {}",
                expectedMontoSinImpuestos, document.getMontoTotalSinImpuestos());
            return false;
        }

        if (!MathUtils.areEqual(document.getIgvTotal(), expectedIgvTotal, TOLERANCE)) {
            logger.error("igvTotal no coincide. Esperado: {}, Obtenido: {}",
                expectedIgvTotal, document.getIgvTotal());
            return false;
        }

        if (!MathUtils.areEqual(document.getMontoTotal(), expectedMontoTotal, TOLERANCE)) {
            logger.error("montoTotal no coincide. Esperado: {}, Obtenido: {}",
                expectedMontoTotal, document.getMontoTotal());
            return false;
        }

        return true;
    }
}
