package com.efact.validator.util;

public final class MathUtils {

    private MathUtils() {
        throw new UnsupportedOperationException();
    }

    public static boolean areEqual(double a, double b, double tolerance) {
        return Math.abs(a - b) < tolerance;
    }
}
