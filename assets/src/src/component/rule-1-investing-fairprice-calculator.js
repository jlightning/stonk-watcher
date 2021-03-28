import {Col, Container, FormControl, InputGroup, Row, Table} from "react-bootstrap";
import {get} from "lodash";
import {humanizeMoney} from "../util/common";
import {useEffect, useState} from "react";

export const Rule1InvestingFairpriceCalculator = ({tickerInfo}) => {
  const [currentEPS, setCurrentEPS] = useState(0);
  const [expectedReturn, setExpectedReturn] = useState(0);
  const [expectedEPSGrowth, setExpectedEPSGrowth] = useState(0);
  const [expectedPE, setExpectedPE] = useState(0);
  const [calculatedFairPrice, setCalculatedFairPrice] = useState(0);

  const eps = get(tickerInfo, 'morningstar_info.financial_data.eps') || get(tickerInfo, 'marketwatch_info.eps') || [];
  const priceOnEarnings = get(tickerInfo, 'morningstar_info.valuation_data.price_on_earnings');
  const equities = get(tickerInfo, 'morningstar_info.financial_data.equities') || get(tickerInfo, 'marketwatch_info.equities') || [];

  useEffect(() => {
    if (eps.length > 0) {
      setCurrentEPS(eps[eps.length - 1].amount);
    }

    setExpectedReturn(get(tickerInfo, 'marketwatch_info.wacc.amount', 0) * 100);
  }, [tickerInfo]);

  useEffect(() => {
    const epsInFuture = currentEPS * Math.pow(1 + expectedEPSGrowth / 100, 10);
    const priceInFuture = epsInFuture * expectedPE;
    const fairprice = priceInFuture / Math.pow(1 + expectedReturn / 100, 10);
    setCalculatedFairPrice(fairprice);
  }, [currentEPS, expectedReturn, expectedEPSGrowth, expectedPE]);

  return (
    <Container>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>EPS Next Year</td>
              <td>{get(tickerInfo, 'finviz_info.eps_next_year.percent')}</td>
              <td>EPS Next 5 Years</td>
              <td>{get(tickerInfo, 'finviz_info.eps_next_5_years.percent')}</td>
            </tr>
          </Table>
        </Col>
      </Row>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>Year</td>
              {
                priceOnEarnings.map(pe => <td>{pe.year.year}</td>)
              }
            </tr>
            <tr>
              <td>P/E</td>
              {
                priceOnEarnings.map(pe => <td>{humanizeMoney(pe.amount.amount)}</td>)
              }
            </tr>
          </Table>
        </Col>
      </Row>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>Year</td>
              {
                equities.map(e => <td>{e.year.year}</td>)
              }
            </tr>
            <tr>
              <td>Equity</td>
              {
                equities.map(e => <td>{humanizeMoney(e.amount)}</td>)
              }
            </tr>
            <tr>
              <td>Equity Growth</td>
              {
                (() => {
                  return equities.map((e, idx) => {
                    if (idx === 0) {
                      return <td>-</td>;
                    }
                    return <td>{`${((e.amount - equities[idx - 1].amount) * 100 / Math.abs(equities[idx - 1].amount)).toFixed(2)}%`}</td>
                  });
                })()
              }
            </tr>
          </Table>
        </Col>
      </Row>
      <Row>
        <Col>
          <Table striped bordered hover style={{fontSize: '11px'}}>
            <tr>
              <td>Year</td>
              {
                eps.map(e => <td>{e.year.year}</td>)
              }
            </tr>
            <tr>
              <td>EPS</td>
              {
                eps.map(e => <td>{humanizeMoney(e.amount)}</td>)
              }
            </tr>
            <tr>
              <td>EPS Growth</td>
              {
                (() => {
                  return eps.map((e, idx) => {
                    if (idx === 0) {
                      return <td>-</td>;
                    }
                    return <td>{`${((e.amount - eps[idx - 1].amount) * 100 / Math.abs(eps[idx - 1].amount)).toFixed(2)}%`}</td>
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
              <InputGroup.Text id="inputGroup-sizing-sm">Current EPS</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={currentEPS}
                         onChange={e => setCurrentEPS(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Est. EPS Growth (In Percentage)</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedEPSGrowth}
                         onChange={e => setExpectedEPSGrowth(parseFloat(e.target.value))}/>
          </InputGroup>
          <InputGroup size="sm" className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="inputGroup-sizing-sm">Expected P/E</InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl aria-label="Small" aria-describedby="inputGroup-sizing-sm" value={expectedPE}
                         onChange={e => setExpectedPE(parseFloat(e.target.value))}/>
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
