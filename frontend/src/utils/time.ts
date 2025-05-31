export const formatCompletionTime = (nanoseconds: number) => {
  const seconds = nanoseconds / 1_000_000_000 // конвертируем наносекунды в секунды
  const days = Math.floor(seconds / (24 * 60 * 60))
  const hours = Math.floor((seconds % (24 * 60 * 60)) / (60 * 60))
  
  if (days > 0) {
    return `${days} д. ${hours} ч.`
  }
  return `${hours} ч.`
} 