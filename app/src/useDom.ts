import React, { useEffect, useRef, useState } from 'react'

export function useElementSize<T extends HTMLElement = HTMLDivElement>(): {
  ref: React.RefObject<T>
  width: number
  height: number
} {
  const { ref, rect } = useBoudingRect<T>()
  return {
    ref,
    width: rect.width,
    height: rect.height,
  }
}

export function useElementWidth<T extends HTMLElement = HTMLDivElement>(): [
  React.RefObject<T>,
  number,
] {
  const { ref, width } = useElementSize<T>()
  return [ref, width]
}

export function useImageWidth(): [
  (event: React.SyntheticEvent<HTMLImageElement, Event>) => void,
  number,
] {
  const [imgWidth, setImgWidth] = useState(0)
  const onImgLoad = (event: React.SyntheticEvent<HTMLImageElement, Event>) => {
    const { naturalWidth } = event.currentTarget
    setImgWidth(naturalWidth)
  }

  return [onImgLoad, imgWidth]
}

export type Rect = {
  width: number
  height: number
  left: number
  top: number
}

export function useBoudingRect<T extends HTMLElement>(): {
  ref: React.RefObject<T>
  rect: Rect
} {
  const [rect, setRect] = useState<Rect>({
    width: 0,
    height: 0,
    left: 0,
    top: 0,
  })
  const ref = useResizeObserver<T>(() => {
    if (ref.current) {
      const rect = ref.current.getBoundingClientRect()
      setRect({
        width: rect.width,
        height: rect.height,
        left: rect.left,
        top: rect.top,
      })
    }
  })

  return { ref, rect }
}

export function useResizeObserver<T extends HTMLElement>(
  callback: (entry: ResizeObserverEntry) => void
): React.RefObject<T> {
  const ref = useRef<T>(null)
  const observer = useRef<ResizeObserver | null>(null)
  useEffect(() => {
    if (ref.current) {
      observer.current = new ResizeObserver((entries) => {
        callback(entries[0])
      })
      observer.current.observe(ref.current)
    }

    return () => {
      if (observer.current) observer.current.disconnect()
    }
  }, [ref])

  return ref
}
