import React, { useRef, useEffect } from 'react'

type Props = {
  width?: number
  height?: number
  points?: [number, number][]
}

export type Plot = {
  x: number
  y: number
  scale: number
  width: number
  height: number
  el: JSX.Element
}

export const Dotplot: React.FC<Props> = ({ width, height }) => {
  const ref = useRef<HTMLCanvasElement>(null)
  const style = {
    background: 'white',
    opacity: 0.7,
    width: '100%',
    height: '100%',
  } as React.CSSProperties

  console.log('dotplot')

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
    ctx.fillStyle = 'rgba(0,255,255,0.5)'

    console.log('draw!')
    for (let i = 0; i < Math.min(width, height); i++) {
      ctx.fillRect(i, i, 1, 1)
    }
  }, [])

  // useEffect(() => {
  //   fetch('http://localhost:8080/')
  //     .then((res)=> res.json())
  //     .then((json) => {
  //       console.log('json', json)
  //     })
  // }, [])
  // const generate = () => {
  //   fetch('http://localhost:8080/generate/')
  // }

  return <canvas ref={ref} width={width} height={height} style={style}></canvas>
}
