import React, { useRef, useEffect } from 'react'

type Props = {
  width?: number
  height?: number
  points?: Points
  dotSize?: number
  colorForward?: string
  colorBackward?: string
}

export type Points = {
  forward: [number, number][]
  backward: [number, number][]
}

export const Dotplot: React.FC<Props> = ({
  width,
  height,
  points = { forward: [], backward: [] },
  dotSize = 1,
  colorForward = '#FF0000',
  colorBackward = '#0000FF',
}) => {
  const ref = useRef<HTMLCanvasElement>(null)
  const style = {
    background: '#FFF',
    opacity: 0.7,
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
    ctx.globalAlpha = 0.5

    console.log('drawing..')
    ctx.fillStyle = colorForward
    for (const point of points.forward) {
      ctx.fillRect(point[0], point[1], dotSize, dotSize)
    }

    ctx.fillStyle = colorBackward
    for (const point of points.backward) {
      ctx.fillRect(point[0], point[1], dotSize, dotSize)
    }

    return () => {
      console.log('removing')
    }
  }, [dotSize])

  return <canvas ref={ref} width={width} height={height} style={style}></canvas>
}
