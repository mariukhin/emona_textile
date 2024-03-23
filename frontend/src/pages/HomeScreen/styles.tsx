import styled from 'styled-components';
import { colors } from 'utils/color';
import { device } from 'utils/deviceBreakpoints';

export const ClientsBlockWrapper = styled.div`
  background-color: ${colors.background.white};
`;

export const ClientsBlockTitleWrapper = styled.div`
  @media ${device.mobile} {
    padding: 80px 0 0;
  }

  @media ${device.laptopL} {
    padding: 80px 0 100px;
  }
`;

export const ClientTicker = styled.div`
  height: 100px;
  width: 100%;
  display: flex;
  flex-direction: row;

  @media ${device.mobile} {
    height: 150px;
  }
  
  @media ${device.laptopL} {
    height: 100px;
  }
`;

export const TickerWrapper = styled.div`
  overflow: hidden;
  width: 100%;
  margin-top: 40px;
  display: block;
  position: relative;
`;

export const LogoContainer = styled.div`
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  align-content: center;
  width: 100%;
  margin: 40px 24px 0; 

  @media ${device.mobile} {
    display: block;
  }

  @media ${device.laptopL} {
    display: none;
  }
`

export const TickerImage = styled.img`
  @media ${device.mobile} {
    vertical-align: middle;
    max-width: 171px;
    margin-right: 30px;
    margin-bottom: 50px;
    max-height: 150px;
  }

  @media ${device.tablet} {
    max-width: 190px;
    margin-right: 50px;
    margin-bottom: 68px;
  }
  

  @media ${device.laptopL} {
    margin-right: 50px;
    margin-bottom: 0;
    max-width: 235px;
    max-height: 100px;

    filter: saturate(0);

    :hover {
      filter: saturate(1);
    }
  }
`;