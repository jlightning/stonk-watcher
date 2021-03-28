import {get} from 'lodash';
import {Col, Container, FormControl, InputGroup, Row, Table} from "react-bootstrap";
import {humanizeMoney} from "../util/common";
import {useEffect, useState} from "react";
import {CanvasJSChart} from "canvasjs-react-charts";

export const DiscountedCashflowFairpriceCalculator = ({tickerInfo}) => {
  const [cashFlows, setCashFlows] = useState([]);
  const [currentCashFlow, setCurrentCashFlow] = useState(0);
  const [expectedReturn, setExpectedReturn] = useState(0);
  const [currentCashFlowGrowth, setCurrentCashFlowGrowth] = useState(0);
  const [calculatedFairPrice, setCalculatedFairPrice] = useState(0);
  const [shareOutstanding, setShareOutstanding] = useState(0);

  useEffect(() => {
    setCashFlows(get(tickerInfo, 'morningstar_info.financial_data.cash_flows') || get(tickerInfo, 'marketwatch_info.free_cash_flow') || []);
    setShareOutstanding(get(tickerInfo, 'finviz_info.share_outstanding', 0));
    setExpectedReturn(get(tickerInfo, 'marketwatch_info.wacc.amount', 0) * 100);
  }, [tickerInfo]);

  useEffect(() => {
    if (cashFlows.length > 0) {
      setCurrentCashFlow(cashFlows[cashFlows.length - 1].amount);
    }
  }, [cashFlows])

  useEffect(() => {
    let arr = [];
    for (let i = 1; i <= 10; i++) {
      arr.push(currentCashFlow * Math.pow(1 + currentCashFlowGrowth / 100, i));
    }

    const perpetualGrowth = 0.025;
    const terminatedCashFlow = arr[arr.length - 1] * (1 + perpetualGrowth) / (expectedReturn/100 - perpetualGrowth);
    arr.push(terminatedCashFlow);

    const discountFactor = expectedReturn/100 + 1;

    arr = arr.map((cf, idx) => {
      return cf / Math.pow(discountFactor, idx === arr.length - 1 ? idx : idx + 1);
    })

    setCalculatedFairPrice((arr.reduce((prev, current) => prev + current, 0)) / shareOutstanding);
  }, [currentCashFlow, currentCashFlowGrowth, expectedReturn]);

  if (!get(tickerInfo, "finviz_info")) {
    return null;
  }

  return (
    <Container>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>Year</td>
              {
                cashFlows.map(cf => <td>{cf.year.year}</td>)
              }
            </tr>
            <tr>
              <td>Cash Flow</td>
              {
                cashFlows.map(cf => <td>{humanizeMoney(cf.amount)}</td>)
              }
            </tr>
            <tr>
              <td>Cash Flow Growth</td>
              {
                (() => {
                  return cashFlows.map((cf, idx) => {
                    if (idx === 0) {
                      return <td>-</td>;
                    }
                    return <td>{`${((cf.amount - cashFlows[idx - 1].amount) * 100 / Math.abs(cashFlows[idx - 1].amount)).toFixed(2)}%`}</td>
                  });
                })()
              }
            </tr>
          </Table>
        </Col>
      </Row>
      <Row>
        <CanvasJSChart options={{
          theme: "light2", // "light1", "dark1", "dark2"
          axisY: {
            title: "Cash Flow",
            suffix: '$',
          },
          axisX: {
            title: "Year",
            prefix: "Y",
            interval: 2
          },
          data: [{
            type: "line",
            toolTipContent: "Year {x}: {y}$",
            dataPoints: cashFlows.map(item => ({
              x: item.year.year,
              y: item.amount,
            }))
          }]
        }}/>
      </Row>
      <Row>
        <Col>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Current Cash Flow</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={currentCashFlow}
                         onChange={e => setCurrentCashFlow(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Est. Cash Flow Growth (In Percentage)</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={currentCashFlowGrowth}
                         onChange={e => setCurrentCashFlowGrowth(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Expected Return (In Percentage)</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedReturn}
                         onChange={e => setExpectedReturn(parseFloat(e.target.value))}/>
          </InputGroup>
        </Col>
        <Col>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Fairprice</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm"
                         value={calculatedFairPrice.toFixed(2)}
                         disabled={true}
                         onChange={e => setCalculatedFairPrice(parseFloat(e.target.value))}/>
          </InputGroup>
        </Col>
      </Row>
    </Container>
  )
}
