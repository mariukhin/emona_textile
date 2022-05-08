// node modules
import React from 'react';
// styles
import {
  BannerContainer,
  BannerWrapper,
  BannerHeading,
  BannerSubHeading,
} from './styles';

type BannerProps = {
  title: string;
  subtitle: string;
}

const Banner: React.FC<BannerProps> = ({ title, subtitle }) => (
  <BannerContainer>
    <BannerWrapper>
      <BannerHeading>{title}</BannerHeading>
      <BannerSubHeading>{subtitle}</BannerSubHeading>
    </BannerWrapper>
  </BannerContainer>
);

export default Banner;
