import { useCallback, useRef } from 'react'

export function useDebounce(callback: () => void, delay: number = 1000) {
  const timer = useRef<number | null>(null)
  const debounce = useCallback(() => {
    if (timer.current) clearTimeout(timer.current)
    timer.current = setTimeout(() => {
      callback()
    }, delay)
  }, [delay, callback])
  return debounce
}
