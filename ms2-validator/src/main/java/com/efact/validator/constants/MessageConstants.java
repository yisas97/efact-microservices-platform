package com.efact.validator.constants;

public final class MessageConstants {

    private MessageConstants() {
        throw new UnsupportedOperationException("Esta es una clase de constantes y no debe ser instanciada");
    }

    public static final String FIRMA_VALIDA = "La firma es válida y el documento no ha sido modificado";
    public static final String FIRMA_INVALIDA = "La firma es inválida o el documento ha sido modificado";
    public static final String DOCUMENTO_NO_ENCONTRADO = "Documento no encontrado";
    public static final String DOCUMENTO_SIN_ITEMS = "El documento no tiene ítems";
    public static final String VALIDACION_ITEMS_FALLIDA = "Validación de ítems fallida";
    public static final String VALIDACION_TOTALES_FALLIDA = "Validación de totales fallida";
}
