import {Col, Container, FormControl, InputGroup, Row, Table} from "react-bootstrap";
import {get} from "lodash";
import {humanizeMoney} from "../util/common";
import {useEffect, useState} from "react";

export const DiscountedDividendFairpriceCalculator = ({tickerInfo}) => {
  const [currentDPS, setCurrentDPS] = useState(0);
  const [expectedReturn, setExpectedReturn] = useState(0);
  const [expectedDPSShortGrowth, setExpectedDPSShortGrowth] = useState(0);
  const [expectedDPSLongGrowth, setExpectedDPSLongGrowth] = useState(0);
  const [expectedDPSShortYears, setExpectedDPSShortYears] = useState(0);
  const [calculatedFairPrice, setCalculatedFairPrice] = useState(0);

  const dividends = get(tickerInfo, 'morningstar_info.dividend_data.dividend_per_shares') || [];

  useEffect(() => {
    if (dividends.length > 0) {
      setCurrentDPS(dividends[dividends.length - 1].amount);
    }
  }, [tickerInfo]);

  useEffect(() => {
    let fairprice = currentDPS * (1 + expectedDPSLongGrowth / 100) / (expectedReturn / 100 - expectedDPSLongGrowth / 100);
    if (expectedDPSShortGrowth > 0 && expectedDPSShortYears > 0) {
      fairprice += currentDPS * expectedDPSShortYears * (expectedDPSShortGrowth / 100 + expectedDPSLongGrowth / 100) / (expectedReturn / 100 - expectedDPSLongGrowth / 100);
    }
    setCalculatedFairPrice(fairprice);
  }, [currentDPS, expectedReturn, expectedDPSShortGrowth, expectedDPSShortYears, expectedDPSLongGrowth]);

  return (
    <Container>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>Year</td>
              {
                dividends.map(d => <td>{d.year.year}</td>)
              }
            </tr>
            <tr>
              <td>Dividend Per Share</td>
              {
                dividends.map(d => <td>{humanizeMoney(d.amount)}</td>)
              }
            </tr>
            <tr>
              <td>Dividend Per Share Growth</td>
              {
                (() => {
                  return dividends.map((e, idx) => {
                    if (idx === 0) {
                      return <td>-</td>;
                    }
                    return <td>{`${((e.amount - dividends[idx - 1].amount) * 100 / Math.abs(dividends[idx - 1].amount)).toFixed(2)}%`}</td>
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
              <InputGroup.Text id="inputGroup-sizing-sm">Current DPS</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={currentDPS}
                         onChange={e => setCurrentDPS(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Est. DPS Growth Short Term(In Percentage)</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedDPSShortGrowth}
                         onChange={e => setExpectedDPSShortGrowth(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Est. DPS Growth Short Term In Year</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedDPSShortYears}
                         onChange={e => setExpectedDPSShortYears(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Est. DPS Growth Long Term(In Percentage)</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedDPSLongGrowth}
                         onChange={e => setExpectedDPSLongGrowth(parseFloat(e.target.value))}/>
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
  );
}
