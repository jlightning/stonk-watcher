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

export const getDebtEquityDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount <= 0.6) {
    return 'good';
  }

  if (amount > 1.1) {
     return 'danger';
  }

  return '';
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

export const getDiscountDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount < 0) {
    return 'danger';
  }

  if (amount > 10 && amount < 20) {
    return 'warn';
  }

  if (amount > 20) {
    return 'good'
  }

  return '';
}

export const getGrossIncomeMarginDangerLevel = (amount) => {
  if (!amount && amount !== 0) {
    return '';
  }

  if (amount <= 0.20) {
    return 'danger';
  }

  if (amount <= 0.40) {
    return 'warn';
  }

  if (amount >= 0.70) {
    return 'good';
  }

  return '';
}

export const humanizeMoney = (amount) => {
  if (Math.abs(amount) > 1000000000000) {
    return `${(amount / 1000000000000).toFixed(2)}T`;
  }
  if (Math.abs(amount) > 1000000000) {
    return `${(amount / 10000000000).toFixed(2)}B`;
  }
  if (Math.abs(amount) > 1000000) {
    return `${(amount / 1000000).toFixed(2)}M`;
  }
  return amount;
}
