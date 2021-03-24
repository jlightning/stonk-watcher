import {
  Button,
  Col,
  Container,
  FormControl,
  InputGroup,
  Modal,
  OverlayTrigger,
  Row,
  Spinner,
  Table
} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css';
// import './bootstrap.min.css';
import './App.css';
import {useEffect, useState} from "react";
import {get, takeRight} from 'lodash';
import {
  getDebtEquityDangerLevel,
  getDiscountDangerLevel,
  getGrossIncomeMarginDangerLevel,
  getPeDangerLevel,
  getReturnColorDangerLevel,
  getRSIDangerLevel,
  getShortFloatDangerLevel,
  shortenString
} from "./util/common";
import {ColorBox} from "./component/colorbox";
import {FairPriceCalculator} from "./component/fairprice-calculator";
import {FaEdit, FaRedo} from "react-icons/all";
import {PerformanceTooltip} from "./component/performance-tooltip";

const SERVER_URL = process.env.REACT_APP_SERVER_URL || '/';

function App() {
  const [tickers, setTickers] = useState([]);
  const [tickerSplatBySector, setTickerSplatBySector] = useState([]);
  const [tickerStr, setTickerStr] = useState('');
  const [details, setDetails] = useState({});
  const [prices, setPrices] = useState({});
  const [show, setShow] = useState(false);
  const [fairPriceTickerInfo, setFairPriceTickerInfo] = useState({});

  useEffect(() => {
    (async () => {
      const result = await fetch(`${SERVER_URL}watchlist`).then(r => r.json());
      setTickers(result);
    })()
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      tickers.forEach(async t => {
        const res = await fetch(`${SERVER_URL}stock/price?ticker=${t}`).then(r => r.json());
        setPrices(prevDate => ({...prevDate, [t]: res}));
      })
    }, 15000);
    return () => {
      clearInterval(interval)
    }
  });

  const loadTickerStockInfo = async t => {
    const res = await fetch(`${SERVER_URL}stock?ticker=${t}`).then(r => r.json());
    setDetails(prevDate => ({...prevDate, [t]: res}));
  }

  useEffect(() => {
    tickers.forEach(loadTickerStockInfo)

    setTickerStr(tickers.reduce((a, b) => a + (a ? ', ' : '') + b, ''));
  }, [tickers]);

  useEffect(() => {
    const arr = tickers.map(t => ({ticker: t, sector: get(details, `['${t}'].finviz_info.sector`, '')}))
      .reduce((prev, cur) => {
        let found = false
        prev.forEach(p => {
          if (p.sector === cur.sector) {
            p.tickers.push(cur.ticker);
            found = true;
          }
        });
        if (!found) {
          prev.push({tickers: [cur.ticker], sector: cur.sector});
        }

        return prev;
      }, []);
    setTickerSplatBySector(arr);
  }, [tickers, details])

  const clearData = async () => {
    await fetch(`${SERVER_URL}stock`, {
      method: 'DELETE',
    })
    setDetails({});
    setTickers(prevState => [...prevState])
  }

  const deleteStockInfo = async (t) => {
    await fetch(`${SERVER_URL}stock?ticker=${t}`, {
      method: 'DELETE',
    })
    setDetails(prevState => ({...prevState, [t]: null}));
    await loadTickerStockInfo(t)
  }

  const updateWatchList = async () => {
    const result = await fetch(`${SERVER_URL}watchlist`, {
      method: 'POST',
      body: tickerStr,
    }).then(r => r.json());
    setTickers(result);
  }

  const showMyFairPriceModal = (t) => {
    setFairPriceTickerInfo(t);
    setShow(true);
  }

  const renderTooltip = (ticker, props, data, title, isPercentage = false) => (
    <PerformanceTooltip key={ticker} props={props} data={data} title={title} isPercentage={isPercentage}/>
  );

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
              <Button variant="danger" onClick={clearData}>Clear data</Button>
            </InputGroup.Append>
          </InputGroup>
        </Row>
        <Row>
          {
            tickers.filter(t => ! get(details, `['${t}']`)).map(t => (
              <Button variant="primary" disabled className='mr-2 mb-2' key={t}>
                <Spinner
                  as="span"
                  animation="border"
                  size="sm"
                  role="status"
                  aria-hidden="true"
                />
                <span className="ml-1">{t}</span>
              </Button>
            ))
          }
        </Row>
        <Row>
          {
            tickerSplatBySector.filter(ts => ts.sector).map((splatTicker, sectorIdx) => (
              <Table striped bordered hover className={'stock-table'} key={splatTicker.sector}>
                <thead>
                <tr>
                  <th><Container><Row className={'row-th'}>Ticker</Row></Container></th>
                  <th><Container><Row className={'row-th'}>Company Name</Row></Container></th>
                  <th><Container><Row className={'row-th'}>RSI</Row></Container></th>
                  <th><Container><Row className={'row-th'}>Shorted</Row></Container></th>
                  <th>
                    <Container>
                      <Row className={'row-th'}>Debt /</Row>
                      <Row className={'row-th'}>Equity</Row>
                    </Container>
                  </th>
                  <th><Container><Row className={'row-th'}>P/B</Row></Container></th>
                  <th>
                    <Container>
                      <Row className={'row-th'}>Gross Income</Row>
                      <Row className={'row-th'}>Margin</Row>
                    </Container>
                  </th>
                  <th>
                    <Container>
                      <Row className={'row-th'}>Net Income</Row>
                      <Row className={'row-th'}>Margin</Row>
                    </Container>
                  </th>
                  <th><Container><Row className={'row-th'}>ROIC</Row></Container></th>
                  <th><Container><Row className={'row-th'}>Sales Growth</Row></Container></th>
                  <th><Container><Row className={'row-th'}>EPS Growth</Row></Container></th>
                  <th><Container><Row className={'row-th'}>Equity Growth</Row></Container></th>
                  <th><Container><Row className={'row-th'}>Cash Flow Growth</Row></Container></th>
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
                <tr>
                  <td colSpan={4}>{splatTicker.sector}</td>
                  <td colSpan={15}>-</td>
                </tr>
                {
                  splatTicker.tickers.map((t, tIdx) => {
                    let rois = [];
                    let roiGrowths = [];
                    let sales = [];
                    let saleGrowths = [];
                    let epsGrowths = [];
                    let eps = [];
                    let equities = [];
                    let equityGrowths = [];
                    let cashFlows = [];
                    let cashFlowGrowths = [];
                    let grossIncomeMargins = [];
                    let netIncomeMargins = [];
                    const detail = get(details, `['${t}']`)
                    if (detail) {
                      grossIncomeMargins = get(detail, 'morningstar_info.financial_data.gross_profit_margins') || get(detail, 'marketwatch_info.gross_income_margin') || [];
                      netIncomeMargins = get(detail, 'morningstar_info.financial_data.net_profit_margins') || get(detail, 'marketwatch_info.net_income_margins') || [];
                      rois = sales = get(detail, 'morningstar_info.rois') || [];
                      roiGrowths = sales = get(detail, 'morningstar_info.roi_growths') || [];
                      sales = get(detail, 'morningstar_info.financial_data.revenues') || get(detail, 'marketwatch_info.sales') || [];
                      saleGrowths = get(detail, 'morningstar_info.financial_data.revenue_growths') || get(detail, 'marketwatch_info.sales_growth') || [];
                      eps = get(detail, 'morningstar_info.financial_data.eps') || get(detail, 'marketwatch_info.eps') || [];
                      epsGrowths = get(detail, 'morningstar_info.financial_data.eps_growths') || get(detail, 'marketwatch_info.eps_growths') || [];
                      equities = get(detail, 'morningstar_info.financial_data.equities') || get(detail, 'marketwatch_info.equities') || [];
                      equityGrowths = get(detail, 'morningstar_info.financial_data.equity_growths') || get(detail, 'marketwatch_info.equity_growths') || [];
                      cashFlows = get(detail, 'morningstar_info.financial_data.cash_flows') || get(detail, 'marketwatch_info.free_cash_flow') || [];
                      cashFlowGrowths = get(detail, 'morningstar_info.financial_data.cash_flow_growths') || get(detail, 'marketwatch_info.free_cash_flow_growths') || [];
                    }
                    let price = get(prices, `['${t}'].price`);
                    price = price || get(details, `['${t}'].finviz_info.price`, '-');
                    let targetPrice = get(details, `['${t}'].finviz_info.target_price`, '-');
                    let msFairPrice = get(details, `['${t}'].morningstar_info.latest_fair_price`, '-');
                    let targetPriceDiscount = (targetPrice - price) * 100 / targetPrice;
                    let msFairPriceDiscount = (msFairPrice - price) * 100 / msFairPrice;

                    return (
                      <>
                        <tr key={t}>
                          <td>{t}</td>
                          <td>{shortenString(get(details, `['${t}'].finviz_info.company_name`, '-'), 12)}</td>
                          <td>
                            <ColorBox
                              dangerLevel={getRSIDangerLevel(get(details, `['${t}'].finviz_info.rsi.amount`))}>{get(details, `['${t}'].finviz_info.rsi.amount`, '-')}</ColorBox>
                          </td>
                          <td>
                            <ColorBox
                              dangerLevel={getShortFloatDangerLevel(get(details, `['${t}'].finviz_info.short_float.amount`))}>{get(details, `['${t}'].finviz_info.short_float.percent`, '-')}</ColorBox>
                          </td>
                          <td>
                            <ColorBox
                              dangerLevel={getDebtEquityDangerLevel(get(details, `['${t}'].finviz_info.debt_on_equity.amount`))}>{get(details, `['${t}'].finviz_info.debt_on_equity.amount`, '-')}</ColorBox>
                          </td>
                          <td>
                            {get(details, `['${t}'].finviz_info.pb.amount`, '-')}
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, grossIncomeMargins, 'Gross Income Margin', true)}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {takeRight(grossIncomeMargins, 1).map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getGrossIncomeMarginDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, netIncomeMargins, 'Net Income Margin', true)}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {takeRight(netIncomeMargins, 1).map(r => (
                                    <Col>{get(r, 'amount.percent', '-')}</Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, rois, 'ROI', true)}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {roiGrowths.map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getReturnColorDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, sales, 'Revenue')}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {saleGrowths.map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getReturnColorDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, eps, 'EPS')}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {epsGrowths.map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getReturnColorDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'right'}
                                            overlay={props => renderTooltip(t, props, equities, 'Equity')}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {equityGrowths.map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getReturnColorDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
                          </td>
                          <td>
                            <OverlayTrigger delay={{show: 50, hide: 150}} placement={'left'}
                                            overlay={props => renderTooltip(t, props, cashFlows, 'Cash Flow')}>
                              <Container className={'can-hover'}>
                                <Row>
                                  {cashFlowGrowths.map(r => (
                                    <Col>
                                      <ColorBox
                                        dangerLevel={getReturnColorDangerLevel(get(r, 'amount.amount'))}>{get(r, 'amount.percent', '-')}</ColorBox>
                                    </Col>
                                  ))}
                                </Row>
                              </Container>
                            </OverlayTrigger>
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
                        <tr>
                          <td colSpan={2}>URL</td>
                          <td colSpan={4}><a href={get(details, `['${t}'].finviz_info.url`)}
                                             target='_blank'>{get(details, `['${t}'].finviz_info.url`)}</a></td>
                          <td colSpan={3}><a href={get(details, `['${t}'].marketwatch_info.url`)}
                                             target='_blank'>{get(details, `['${t}'].marketwatch_info.url`)}</a></td>
                          <td colSpan={4}><a href={get(details, `['${t}'].morningstar_info.url`)}
                                             target='_blank'>{get(details, `['${t}'].morningstar_info.url`)}</a></td>
                          <td colSpan={3}>-</td>
                          <td colSpan={2}>Action</td>
                          <td>
                            <Container className={'action-button'}>
                              <Row>
                                <Col>
                                  <button className={'no-style'} onClick={e => showMyFairPriceModal(detail)}><FaEdit/>
                                  </button>
                                </Col>
                                <Col>
                                  <button className={'no-style'} onClick={e => deleteStockInfo(t)}><FaRedo/>
                                  </button>
                                </Col>
                              </Row>
                            </Container>
                          </td>
                        </tr>
                      </>
                    )
                  })
                }
                </tbody>
              </Table>
            ))
          }
        </Row>
      </Container>

      <Modal
        show={show}
        onHide={e => setShow(false)}
        backdrop="static"
        keyboard={false}
        className={'fairprice-edit-modal'}
      >
        <Modal.Header closeButton>
          <Modal.Title>Fair Price Calculator for {get(fairPriceTickerInfo, 'finviz_info.company_name')}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <FairPriceCalculator tickerInfo={fairPriceTickerInfo}/>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary">Understood</Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}

export default App;
