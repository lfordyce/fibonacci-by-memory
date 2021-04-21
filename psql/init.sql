CREATE TABLE IF NOT EXISTS fib_store(
    ordinal   INTEGER,
    fibonacci BIGINT NOT NULL,
    PRIMARY KEY (ordinal)
);

INSERT INTO fib_store (ordinal, fibonacci) VALUES (0, 0);
INSERT INTO fib_store (ordinal, fibonacci) VALUES (1, 1);

-- fibonacci memoization function with the fib_store table as the cache
CREATE OR REPLACE FUNCTION fibonacci_cached(ordinal_for INTEGER)
    RETURNS BIGINT AS
$$
DECLARE
    ret BIGINT;
BEGIN
    IF ordinal_for < 2 THEN
        RETURN ordinal_for;
    END IF;

    SELECT INTO ret fibonacci FROM fib_store WHERE ordinal = ordinal_for;

    IF ret IS NULL THEN
        ret := fibonacci_cached(ordinal_for - 2) + fibonacci_cached(ordinal_for - 1);
        INSERT INTO fib_store(ordinal, fibonacci)
        VALUES (ordinal_for, ret);
    END IF;
    RETURN ret;
END;
$$ LANGUAGE plpgsql;

-- FUNCTION
CREATE OR REPLACE FUNCTION base_case_fibonacci()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO fib_store (ordinal, fibonacci) VALUES (0, 0);
    INSERT INTO fib_store (ordinal, fibonacci) VALUES (1, 1);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER
CREATE TRIGGER trigger_base_case_insert
    AFTER TRUNCATE
    ON fib_store
EXECUTE PROCEDURE base_case_fibonacci();