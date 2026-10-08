type Row = Record<string, string | number | boolean | null>;

interface Props {
  atributes: Array<string>;
  list: Array<Row>;
}

function TableX({ atributes, list }: Props) {
  return (
    <table className="table">
      <thead>
        <tr>
          {atributes.map((a) => {
            return <th scope="col">{a}</th>;
          })}
        </tr>
      </thead>
      <tbody>
        {list.map((ins) => {
          return (
            <tr>
              {atributes.map((ca) => {
                return <td>{ins[ca]}</td>;
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

export default TableX;

/*

<table class="table">
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">First</th>
      <th scope="col">Last</th>
      <th scope="col">Handle</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th scope="row">1</th>
      <td>Mark</td>
      <td>Otto</td>
      <td>@mdo</td>
    </tr>
    <tr>
      <th scope="row">2</th>
      <td>Jacob</td>
      <td>Thornton</td>
      <td>@fat</td>
    </tr>
    <tr>
      <th scope="row">3</th>
      <td colspan="2">Larry the Bird</td>
      <td>@twitter</td>
    </tr>
  </tbody>
</table>


*/
