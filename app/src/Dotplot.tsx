import React, { useRef, useEffect } from 'react'

type Props = {
  width?: number
  height?: number
  points?: [number, number][]
}

export const Dotplot: React.FC<Props> = ({ width, height, points = [] }) => {
  const ref = useRef<HTMLCanvasElement>(null)
  const style = {
    background: 'white',
    // opacity: 0.7,
    width: '100%',
    height: '100%',
    imageRendering: 'pixelated',
    // image-rendering: crisp-edges;
  } as React.CSSProperties

  console.log('dotplot', width, height)

  // draw wave to canvas
  useEffect(() => {
    if (!ref.current) return

    const canvas = ref.current
    const ctx = canvas.getContext('2d')
    const width = canvas.width
    const height = canvas.height
    if (!ctx) return

    // reset canvas
    ctx.clearRect(0, 0, width, height)
    // ctx.fillStyle = 'rgba(0,255,255,0.5)'
    ctx.fillStyle = 'red'

    console.log('drawing..')
    for (const point of points) {
      ctx.fillRect(point[0], point[1], 1, 1)
    }

    return () => {
      console.log('removing')
    }
  }, [])

  return <canvas ref={ref} width={width} height={height} style={style}></canvas>
}
