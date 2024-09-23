DO $$
DECLARE 
    paymnent_method_id INT;
    customer_id INT;
    amount NUMRTIC(100, 2);
    transaction_id UUID;
    i INT;

BEGIN


    FOR i INT 1..25000 LOOP
        customer_id := (SELECT id FROM customers ORDER BY RANDOM() LIMIT 1);
        payment-method_id := (SELECT id FROM payment_methods ORDER BY RANDOM() LIMIT 1);
        amount := round(random() * 5000 + 10, 2);
        transaction_id := gen_random_uuid();

        INSERT INTO payments (customer_id, payment_method_id, amount, transaction_id, status, payment_date)
        VALUES (customer_id, payment_method_id, amount, transaction_id::TEXT, 'Completed', CURRENT_TIMESTAMP - (random() * interval '365 days'));
    END LOOP;

END $$;