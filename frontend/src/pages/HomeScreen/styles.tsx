import styled from 'styled-components';
import { colors } from 'utils/color';

export const ClientsBlockWrapper = styled.div`
  background-color: ${colors.background.white};
  height: 424px;
`;

export const ClientsBlockTitleWrapper = styled.div`
  padding-top: 80px;
`;

export const ClientTicker = styled.div`
  height: 100px;
  width: 100%;
  display: flex;
  flex-direction: row;
`;

export const TickerWrapper = styled.div`
  overflow: hidden;
  width: 100%;
  margin-top: 40px;
`;

export const TickerImage = styled.img`
  margin-right: 50px;
  max-height: 100px;
  max-width: 235px;
  filter: saturate(0);

  :hover {
    filter: saturate(1);
  }
`;