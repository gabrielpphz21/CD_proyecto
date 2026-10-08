import { useLayoutEffect, useRef, type RefObject } from "react";

interface Props {
  image_path: string;
  xy_s: number[][];
}

function transformCoordinates(
  x: number,
  y: number,
  w: number,
  h: number,
): number[] {
  var new_x = x + w / 2;
  var new_y = -y + h / 2;
  return [new_x, new_y];
}

function renderCanvas(
  canvasRef: RefObject<HTMLCanvasElement | null>,
  image_path: string,
  xy_s: number[][],
) {
  const img = new Image(50, 50); // Optional width and height
  img.onload = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    var ctx = canvas.getContext("2d");
    if (!ctx) {
      return;
    }
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    for (const [x, y] of xy_s) {
      ctx.drawImage(img, x, y, 25, 25);
    }
  };
  img.onerror = () => {
    console.error("Failed to load:", image_path);
  };
  img.src = image_path;
}

function CanvasMap({ image_path, xy_s }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const transformedCoordinates = xy_s.map(([x, y]) =>
    transformCoordinates(x, y, 800, 800),
  );
  console.log("image path:", image_path);
  useLayoutEffect(() => {
    renderCanvas(canvasRef, image_path, transformedCoordinates);
  }, [image_path, xy_s]);
  return <canvas ref={canvasRef} width={800} height={800} />;
}

export default CanvasMap;
