------------------function to add customers automatically------------------------
CREATE OR REPLACE FUNCTION add_customer()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO customers (full_name, email, phone_number)
    VALUES ('Customer_' || NEW.id, 'customer_' || NEW.id || '@example.com', '123-456-7890');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

---------------------Create a trigger to add dat----------------------------------------------
CEWATE TRIGGER add_customer_trigger
ALTER INSERT ON payments
FOR EACH ROW
EXECUTE FUNCTION add_customer();
