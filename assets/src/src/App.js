import './App.css';
import {Button, Col, Container, Row, Table} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css';
import {useEffect, useState} from "react";
import {get} from 'lodash';
import {getPeDangerLevel, getReturnColorDangerLevel, getRSIDangerLevel, getShortFloatDangerLevel} from "./util/common";
import {ColorBox} from "./component/colorbox";

const SERVER_URL = 'http://localhost:8080/'

function App() {
  const [tickers, setTickers] = useState([]);
  const [details, setDetails] = useState({})

  useEffect(() => {
    (async () => {
      const result = await fetch(`${SERVER_URL}watchlist`).then(r => r.json());
      setTickers(result);
    })()
  }, [])

  useEffect(() => {
    tickers.forEach(async t => {
      const res = await fetch(`${SERVER_URL}stock?ticker=${t}`).then(r => r.json());
      setDetails(prevDate => ({...prevDate, [t]: res}));
    })
  }, [tickers])

  const clearData = async () => {
    await fetch(`${SERVER_URL}stock`, {
      method: 'DELETE',
    })
    setDetails({});
    setTickers(prevState => [...prevState])
  }

  return (
    <div className="App">
      <Table striped bordered hover className={'stock-table'}>
        <thead>
        <tr>
          <th>Ticker</th>
          <th>Company Name</th>
          <th>RSI</th>
          <th>Shorted</th>
          <th>ROIC (10 5 1)</th>
          <th>Sales Growth (5 3 1)</th>
          <th>EPS Growth (5 3 1)</th>
          <th>Equity Growth (5 3 1)</th>
          <th>Cash Flow Growth (5 3 1)</th>
          <th>EPS (TTM)</th>
          <th>Current P/E</th>
          <th>Price</th>
          <th>Target Price</th>
          <th>MS Fair Price</th>
        </tr>
        </thead>
        <tbody>
        {
          tickers.map(t => {
            let roiGrowths = [];
            let saleGrowths = [];
            let epsGrowths = [];
            let equityGrowths = [];
            let cashFlowGrowths = [];
            const detail = get(details, `['${t}']`)
            if (detail) {
              roiGrowths = [
                get(detail, 'morningstar_info.roi_10_years', '-'),
                get(detail, 'morningstar_info.roi_5_years', '-'),
                get(detail, 'morningstar_info.roi_last_year', '-'),
              ];
              saleGrowths = [
                get(detail, `marketwatch_info.sales_growth_5_years`, '-'),
                get(detail, `marketwatch_info.sales_growth_3_years`, '-'),
                get(detail, `marketwatch_info.sales_growth_last_year`, '-'),
              ];
              epsGrowths = [
                get(detail, `marketwatch_info.eps_growth_5_years`, '-'),
                get(detail, `marketwatch_info.eps_growth_3_years`, '-'),
                get(detail, `marketwatch_info.eps_growth_last_year`, '-'),
              ];
              equityGrowths = [
                get(detail, `marketwatch_info.equity_growth_5_years`, '-'),
                get(detail, `marketwatch_info.equity_growth_3_years`, '-'),
                get(detail, `marketwatch_info.equity_growth_last_year`, '-'),
              ]
              cashFlowGrowths = [
                get(detail, `marketwatch_info.free_cash_flow_growth_5_years`, '-'),
                get(detail, `marketwatch_info.free_cash_flow_growth_3_years`, '-'),
                get(detail, `marketwatch_info.free_cash_flow_growth_last_year`, '-'),
              ]
            }
            return (
              <tr>
                <td>{t}</td>
                <td>{get(details, `['${t}'].finviz_info.company_name`, '-')}</td>
                <td>
                  <ColorBox
                    dangerLevel={getRSIDangerLevel(get(details, `['${t}'].finviz_info.rsi.amount`))}>{get(details, `['${t}'].finviz_info.rsi.amount`, '-')}</ColorBox>
                </td>
                <td>
                  <ColorBox
                    dangerLevel={getShortFloatDangerLevel(get(details, `['${t}'].finviz_info.short_float.amount`))}>{get(details, `['${t}'].finviz_info.short_float.percent`, '-')}</ColorBox>
                </td>
                <td>
                  <Container>
                    <Row>
                      {roiGrowths.map(r => (
                        <Col>
                          <ColorBox
                            dangerLevel={getReturnColorDangerLevel(get(r, 'amount'))}>{get(r, 'percent', '-')}</ColorBox>
                        </Col>
                      ))}
                    </Row>
                  </Container>
                </td>
                <td>
                  <Container>
                    <Row>
                      {saleGrowths.map(r => (
                        <Col>
                          <ColorBox
                            dangerLevel={getReturnColorDangerLevel(get(r, 'amount'))}>{get(r, 'percent', '-')}</ColorBox>
                        </Col>
                      ))}
                    </Row>
                  </Container>
                </td>
                <td>
                  <Container>
                    <Row>
                      {epsGrowths.map(r => (
                        <Col>
                          <ColorBox
                            dangerLevel={getReturnColorDangerLevel(get(r, 'amount'))}>{get(r, 'percent', '-')}</ColorBox>
                        </Col>
                      ))}
                    </Row>
                  </Container>
                </td>
                <td>
                  <Container>
                    <Row>
                      {equityGrowths.map(r => (
                        <Col>
                          <ColorBox
                            dangerLevel={getReturnColorDangerLevel(get(r, 'amount'))}>{get(r, 'percent', '-')}</ColorBox>
                        </Col>
                      ))}
                    </Row>
                  </Container>
                </td>
                <td>
                  <Container>
                    <Row>
                      {cashFlowGrowths.map(r => (
                        <Col>
                          <ColorBox
                            dangerLevel={getReturnColorDangerLevel(get(r, 'amount'))}>{get(r, 'percent', '-')}</ColorBox>
                        </Col>
                      ))}
                    </Row>
                  </Container>
                </td>
                <td>
                  {get(details, `['${t}'].finviz_info.epsttm`, '-')}
                </td>
                <td>
                  <ColorBox
                    dangerLevel={getPeDangerLevel(get(details, `['${t}'].finviz_info.pe.amount`))}>{get(details, `['${t}'].finviz_info.pe.amount`, '-')}</ColorBox>
                </td>
                <td>{get(details, `['${t}'].finviz_info.price`, '-')}</td>
                <td>{get(details, `['${t}'].finviz_info.target_price`, '-')}</td>
                <td>{get(details, `['${t}'].morningstar_info.latest_fair_price`, '-')}</td>
              </tr>
            )
          })
        }
        </tbody>
      </Table>
      <Button variant="danger" onClick={clearData}>Clear data</Button>
    </div>
  );
}

export default App;
