CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,
    report_name TEXT NOT NULL,
    client_email TEXT NOT NULL,
    message_received_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE investment_results (
    client_email TEXT NOT NULL,
    report_name TEXT NOT NULL,
    gains DECIMAL NOT NULL,
    losses DECIMAL NOT NULL,
    info_date TIMESTAMPTZ
);

/*

{
  "clientName": "Acme Corp",
  "reportName": "Monthly_Sales_Sept",
  "email": "jane.doe@acme.com"
}

*/

--docker build -t mypg .   
--docker run -d --name database -p 5432:5432 mypg 