export const hoursToString = (number: number): string =>
  Math.floor(number) + ':' + Math.round((number - Math.floor(number)) * 60).toString().padStart(2, '0')

export const formatMoney = (number: number): string => {
  return `$` + parseInt(number.toString())
}