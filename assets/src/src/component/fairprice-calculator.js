import {Col, Container, Dropdown, DropdownButton, Row} from "react-bootstrap";
import {useState} from "react";
import {find, get} from "lodash";
import {DiscountedCashflowFairpriceCalculator} from "./discounted-cashflow-fairprice-calculator";
import {Rule1InvestingFairpriceCalculator} from "./rule-1-investing-fairprice-calculator";

export const FairpriceCalculator = ({tickerInfo}) => {
  const [value, setValue] = useState('dcf');

  const options = [
    {title: 'Discounted Cash Flow', key: 'dcf'},
    {title: 'Rule #1 investing', key: 'rule1'},
  ];

  return (
    <>
      <Container className='mb-3'>
        <Row>
          <Col>
            <DropdownButton
              alignRight
              title={get(find(options, o => o.key === value), 'title', 'Please select')}
              id="dropdown-menu-align-right"
              onSelect={e => setValue(e)}
            >
              {
                options.map(o => <Dropdown.Item eventKey={o.key}>{o.title}</Dropdown.Item>)
              }
            </DropdownButton>
          </Col>
        </Row>
      </Container>
      {
        value === 'dcf' && <DiscountedCashflowFairpriceCalculator tickerInfo={tickerInfo}/>
      }
      {
        value === 'rule1' && <Rule1InvestingFairpriceCalculator tickerInfo={tickerInfo}/>
      }
    </>
  )
}
