import React, { useRef, useEffect } from 'react'

type Props = {
  width?: number
  height?: number
  points?: Points
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
  colorForward = '#FF0000',
  colorBackward = '#0000FF',
}) => {
  const ref = useRef<HTMLCanvasElement>(null)
  const style = {
    background: 'white',
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

    console.log('drawing..')
    ctx.fillStyle = colorForward
    for (const point of points.forward) {
      ctx.fillRect(point[0], point[1], 1, 1)
    }

    ctx.fillStyle = colorBackward
    for (const point of points.backward) {
      ctx.fillRect(point[0], point[1], 1, 1)
    }

    return () => {
      console.log('removing')
    }
  }, [])

  return <canvas ref={ref} width={width} height={height} style={style}></canvas>
}
