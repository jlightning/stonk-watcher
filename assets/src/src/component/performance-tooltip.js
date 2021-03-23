import {Popover, Table} from "react-bootstrap";
import {get} from "lodash";
import {getReturnColorDangerLevel, humanizeMoney} from "../util/common";
import {CanvasJSChart} from "canvasjs-react-charts";
import {ColorBox} from "./colorbox";

export const PerformanceTooltip = ({props, isPercentage, title, data}) => {
  return (
    <Popover id="button-tooltip" {...props}>
      <Popover.Title>{title}</Popover.Title>
      <Popover.Content>
        <Table striped bordered hover>
          <tr>
            <td>Year</td>
            {data.map(item => (
              <td>{get(item, 'year.year')}</td>
            ))}
          </tr>
          {
            isPercentage ? (
              <>
                <tr>
                  <td>Amount</td>
                  {data.map(item => (
                    <td>{get(item, 'amount.percent')}</td>
                  ))}
                </tr>
              </>
            ) : (
              <>
                <tr>
                  <td>Amount</td>
                  {data.map(item => (
                    <td>{humanizeMoney(get(item, 'amount'))}</td>
                  ))}
                </tr>
                <tr>
                  <td>Percentage</td>
                  {
                    data.map((item, idx) => {
                      if (idx === 0) {
                        return <td>-</td>;
                      }
                      if (!data[idx - 1].amount) {
                        return <td>-</td>;
                      }
                      const percentage = ((item.amount - data[idx - 1].amount) * 100 / Math.abs(data[idx - 1].amount)).toFixed(2);
                      return <td><ColorBox dangerLevel={getReturnColorDangerLevel(percentage/100)}>{percentage}%</ColorBox></td>;
                    })
                  }
                </tr>
              </>
            )
          }
        </Table>
        <CanvasJSChart options={{
          theme: "light2", // "light1", "dark1", "dark2"
          axisY: {
            title: title,
            suffix: isPercentage ? '%' : '$',
          },
          axisX: {
            title: "Year",
            prefix: "Y",
            interval: 2
          },
          data: [{
            type: "line",
            toolTipContent: "Year {x}: {y}%",
            dataPoints: data.map(item => ({
              x: item.year.year,
              y: isPercentage ? item.amount.amount * 100 : item.amount,
            }))
          }]
        }}/>
      </Popover.Content>
    </Popover>
  )
}
