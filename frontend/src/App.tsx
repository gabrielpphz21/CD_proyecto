import "./App.css";
import satelite from "./assets/satellite.png";
import CanvasMap from "./components/CanvasMap";
import Navbar from "./components/Navbar";
import SateliteView from "./components/SateliteView";
import RocketView from "./components/RocketView";
import EnergyTransactionView from "./components/EnergyTransactionView";
import TenantView from "./components/TenantView";
import DestinationView from "./components/DestinationView";
import TakeOffView from "./components/TakeOffView";
import SateliteTable from "./components/SateliteTable";
import RocketTable from "./components/RocketsTable";
import EnergyTransactionsTable from "./components/EnergyTransactionsTable";
import TenantTable from "./components/TenantsTable";
import TakeOffTable from "./components/TakeOffTable";
import DestinationTable from "./components/TakeOffTable";

function App() {
  var coo = [
    [-200, 200],
    [200, 200],
    [200, -200],
    [-200, -200],
    [-100, -200],
  ];

  return (
    <div>
      <Navbar></Navbar>
      <CanvasMap image_path={satelite} xy_s={coo} />
      <SateliteView satelite_id={2}></SateliteView>
      <RocketView rocket_id={1}></RocketView>
      <EnergyTransactionView energyTransaction_id={1}></EnergyTransactionView>
      <TenantView tenant_id={2}></TenantView>
      <DestinationView destination_id={1}></DestinationView>
      <TakeOffView takeOff_id={1}></TakeOffView>
      <SateliteTable></SateliteTable>
      <RocketTable></RocketTable>
      <EnergyTransactionsTable></EnergyTransactionsTable>
      <TenantTable></TenantTable>
      <TakeOffTable></TakeOffTable>
      <DestinationTable></DestinationTable>
    </div>
  );
}

export default App;
