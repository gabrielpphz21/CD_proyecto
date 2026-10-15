import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface Tenant extends Row {
  id: number;
  name: string;
  tax_id: string;
}

function TenantTable() {
  const [tenants, setTenants] = useState<Tenant[]>();

  useEffect(() => {
    async function loadTenants() {
      const response = await fetch(`http://localhost:8000/tenant/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setTenants(data);
    }

    loadTenants();
  }, []);

  return (
    <div>
      <h1>TENANTS</h1>
      {tenants && (
        <TableX atributes={Object.keys(tenants[0])} list={tenants}></TableX>
      )}
    </div>
  );
}

export default TenantTable;
