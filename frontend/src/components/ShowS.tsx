interface Props {
  s_name: string;
  s_data: Record<string, unknown>;
}

function renderObject(
  t: string,
  v: Record<string, unknown>,
): React.ReactNode[] {
  const toRender: React.ReactNode[] = [];

  toRender.push(<h1 key={`title-${t}`}>{t}</h1>);

  for (const [key, val] of Object.entries(v)) {
    if (val !== null && typeof val === "object" && !Array.isArray(val)) {
      toRender.push(...subRender(key, val as Record<string, unknown>));
    } else {
      toRender.push(
        <p key={key}>
          <span className="badge bg-primary">{key}</span> : {String(val)}
        </p>,
      );
    }
  }

  return toRender;
}

function subRender(t: string, v: Record<string, unknown>): React.ReactNode[] {
  const toRender: React.ReactNode[] = [];

  toRender.push(<h3 key={`title-${t}`}>{t}</h3>);

  for (const [key, val] of Object.entries(v)) {
    if (val !== null && typeof val === "object" && !Array.isArray(val)) {
      toRender.push(...subRender(key, val as Record<string, unknown>));
    } else {
      toRender.push(
        <p key={key}>
          <span className="badge bg-primary">{key}</span> : {String(val)}
        </p>,
      );
    }
  }

  return toRender;
}

function ShowS({ s_name, s_data }: Props) {
  return <div>{renderObject(s_name, s_data)}</div>;
}

export default ShowS;
