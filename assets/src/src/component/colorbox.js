const style = {
  paddingLeft: '3px',
  paddingRight: '3px',
}

const goodColorClass = 'bg-success';
const warningColorClass = 'bg-warning';
const dangerColorClass = 'bg-danger';

export const ColorBox = ({children, dangerLevel}) => {
  let className = '';
  switch (dangerLevel) {
    case 'good':
      className = goodColorClass;
      break;
    case 'warn':
      className = warningColorClass;
      break;
    case 'danger':
      className = dangerColorClass;
  }
  return (
    <div style={style} className={className}>
      {children}
    </div>
  );
}
