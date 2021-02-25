export const getReturnColorDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount <= 0) {
    return 'danger'
  }
  if (amount < 0.1) {
    return 'warn';
  }

  return 'good';
}

export const getRSIDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount <= 40) {
    return 'good'
  }

  if (amount >= 80) {
    return 'danger'
  }

  if (amount >= 60) {
    return 'warn'
  }

  return '';
}

export const getShortFloatDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount <= 0.06) {
    return 'good';
  }

  if (amount <= 0.08) {
    return 'warn';
  }

  return 'danger';
}

export const getPeDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return 'danger';
  }

  if (amount > 100) {
    return 'danger'
  }

  return ''
}
