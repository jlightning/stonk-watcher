import {get} from 'lodash';
import {Col, Container, FormControl, InputGroup, Row, Table} from "react-bootstrap";
import {humanizeMoney} from "../util/common";
import {useEffect, useState} from "react";

export const FairPriceCalculator = ({tickerInfo}) => {
  const [currentCashFlow, setCurrentCashFlow] = useState(0);
  const [expectedReturn, setExpectedReturn] = useState(0);
  const [currentCashFlowGrowth, setCurrentCashFlowGrowth] = useState(0);
  const [calculatedFairPrice, setCalculatedFairPrice] = useState(0);
  const [shareOutstanding, setShareOutstanding] = useState(0);

  useEffect(() => {
    const arr = get(tickerInfo, 'marketwatch_info.free_cash_flow', []);
    if (arr.length > 0) {
      setCurrentCashFlow(arr[arr.length - 1]);
    }

    setShareOutstanding(get(tickerInfo, 'finviz_info.share_outstanding', 0));
  }, [tickerInfo]);

  useEffect(() => {
    let arr = [];
    for (let i = 1; i <= 10; i++) {
      arr.push(currentCashFlow * Math.pow((1 + currentCashFlowGrowth / 100), i));
    }

    arr = arr.map(cf => cf - cf * expectedReturn / 100);
    setCalculatedFairPrice((arr.reduce((prev, current) => prev + current, 0) + arr[arr.length - 1] * (currentCashFlowGrowth - expectedReturn)) / shareOutstanding);
  }, [currentCashFlow, currentCashFlowGrowth, expectedReturn]);

  if (!get(tickerInfo, "finviz_info")) {
    return null;
  }

  return (
    <Container>
      <Row>
        <Col>
          <Table striped bordered hover>
            <tr>
              <td>Year</td>
              {
                get(tickerInfo, 'marketwatch_info.years', []).map(y => <td>{y}</td>)
              }
            </tr>
            <tr>
              <td>Cash Flow</td>
              {
                get(tickerInfo, 'marketwatch_info.free_cash_flow', []).map(cf => <td>{humanizeMoney(cf)}</td>)
              }
            </tr>
            <tr>
              <td>Cash Flow Growth</td>
              {
                (() => {
                  const arr = get(tickerInfo, 'marketwatch_info.free_cash_flow', []);
                  return arr.map((cf, idx) => {
                    if (idx === 0) {
                      return <td>-</td>;
                    }
                    return <td>{`${((cf - arr[idx - 1]) * 100 / Math.abs(arr[idx - 1])).toFixed(2)}%`}</td>
                  });
                })()
              }
            </tr>
          </Table>
        </Col>
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
