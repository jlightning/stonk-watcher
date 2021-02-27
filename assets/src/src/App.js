import './App.css';
import {Button, Col, Container, FormControl, InputGroup, Row, Table} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css';
import {useEffect, useState} from "react";
import {get} from 'lodash';
import {
  getDiscountDangerLevel,
  getPeDangerLevel,
  getReturnColorDangerLevel,
  getRSIDangerLevel,
  getShortFloatDangerLevel
} from "./util/common";
import {ColorBox} from "./component/colorbox";

const SERVER_URL = 'http://localhost:8080/'

function App() {
  const [tickers, setTickers] = useState([]);
  const [tickerStr, setTickerStr] = useState('');
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

    setTickerStr(tickers.reduce((a, b) => a + (a ? ', ' : '') + b, ''));
  }, [tickers])

  const clearData = async () => {
    await fetch(`${SERVER_URL}stock`, {
      method: 'DELETE',
    })
    setDetails({});
    setTickers(prevState => [...prevState])
  }

  const updateWatchList = async () => {
    const result =await fetch(`${SERVER_URL}watchlist`, {
      method: 'POST',
      body: tickerStr,
    }).then(r => r.json());
    setTickers(result);
  }

  return (
    <div className="App">
      <Container className={'App-container'}>
        <Row>
          <InputGroup className="mb-3">
            <InputGroup.Prepend>
              <InputGroup.Text id="basic-addon3">
                Watchlist
              </InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl
              id="basic-url"
              aria-describedby="basic-addon3"
              type='text'
              value={tickerStr}
              onChange={e => setTickerStr(e.target.value)}
            />
            <InputGroup.Append>
              <Button variant="primary" onClick={updateWatchList}>Update</Button>
            </InputGroup.Append>
          </InputGroup>
        </Row>
        <Row>
          <Table striped bordered hover className={'stock-table'}>
            <thead>
            <tr>
              <th><Container><Row className={'row-th'}>Ticker</Row></Container></th>
              <th><Container><Row className={'row-th'}>Company Name</Row></Container></th>
              <th><Container><Row className={'row-th'}>RSI</Row></Container></th>
              <th><Container><Row className={'row-th'}>Shorted</Row></Container></th>
              <th><Container><Row className={'row-th'}>ROIC (10 5 1)</Row></Container></th>
              <th><Container><Row className={'row-th'}>Sales Growth (5 3 1)</Row></Container></th>
              <th><Container><Row className={'row-th'}>EPS Growth (5 3 1)</Row></Container></th>
              <th><Container><Row className={'row-th'}>Equity Growth (5 3 1)</Row></Container></th>
              <th><Container><Row className={'row-th'}>Cash Flow Growth (5 3 1)</Row></Container></th>
              <th>
                <Container>
                  <Row className={'row-th'}>Dividend</Row>
                  <Row className={'row-th'}>Yield</Row>
                </Container>
              </th>
              <th><Container><Row className={'row-th'}>EPS (TTM)</Row></Container></th>
              <th>
                <Container>
                  <Row className={'row-th'}>Current</Row>
                  <Row className={'row-th'}>P/E</Row>
                </Container>
              </th>
              <th><Container><Row className={'row-th'}>Price</Row></Container></th>
              <th>
                <Container>
                  <Row className={'row-th'}>Target</Row>
                  <Row className={'row-th'}>Price</Row>
                </Container>
              </th>
              <th>
                <Container>
                  <Row className={'row-th'}>MS</Row>
                  <Row className={'row-th'}>Fair Price</Row>
                </Container>
              </th>
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
                let price = get(details, `['${t}'].finviz_info.price`, '-');
                let targetPrice = get(details, `['${t}'].finviz_info.target_price`, '-');
                let msFairPrice = get(details, `['${t}'].morningstar_info.latest_fair_price`, '-');
                let targetPriceDiscount = (targetPrice - price) * 100 / price;
                let msFairPriceDiscount = (msFairPrice - price) * 100 / price;
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
                      {get(details, `['${t}'].finviz_info.dividend_yield.percent`, '-')}
                    </td>
                    <td>
                      {get(details, `['${t}'].finviz_info.epsttm`, '-')}
                    </td>
                    <td>
                      <ColorBox
                        dangerLevel={getPeDangerLevel(get(details, `['${t}'].finviz_info.pe.amount`))}>{get(details, `['${t}'].finviz_info.pe.amount`, '-')}</ColorBox>
                    </td>
                    <td>{price}</td>
                    <td>
                      <Container>
                        <Row>
                          <Col>{targetPrice}</Col>
                          <Col><ColorBox
                            dangerLevel={getDiscountDangerLevel(targetPriceDiscount)}>{`${targetPriceDiscount.toFixed(2)}%`}</ColorBox></Col>
                        </Row>
                      </Container>
                    </td>
                    <td>
                      <Container>
                        <Row>
                          <Col>{msFairPrice}</Col>
                          <Col><ColorBox
                            dangerLevel={getDiscountDangerLevel(msFairPriceDiscount)}>{`${msFairPriceDiscount.toFixed(2)}%`}</ColorBox></Col>
                        </Row>
                      </Container>
                    </td>
                  </tr>
                )
              })
            }
            </tbody>
          </Table>
          <Button variant="danger" onClick={clearData}>Clear data</Button>
        </Row>
      </Container>
    </div>
  );
}

export default App;
