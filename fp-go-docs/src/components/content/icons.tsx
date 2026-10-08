import React from 'react';

type IconProps = {size?: number; className?: string};

export const InfoIcon = ({size = 18, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M16 2A14 14 0 1 0 30 16 14 14 0 0 0 16 2zm0 26a12 12 0 1 1 12-12 12 12 0 0 1-12 12z" />
    <path d="M16 14h-2v2h1v6h-1v2h4v-2h-1v-8zm-.5-5a1.5 1.5 0 1 0 1.5 1.5A1.5 1.5 0 0 0 15.5 9z" />
  </svg>
);

export const WarnIcon = ({size = 18, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M16 23a1.5 1.5 0 1 0 1.5 1.5A1.5 1.5 0 0 0 16 23zM15 13h2v8h-2z" />
    <path d="M28.7 26.31L17.36 4.27a1.51 1.51 0 0 0-2.72 0L3.3 26.31a1.5 1.5 0 0 0 1.36 2.19h22.69a1.5 1.5 0 0 0 1.35-2.19z" />
  </svg>
);

export const CheckIcon = ({size = 14, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M13 24l-9-9 1.4-1.4L13 21.2 26.6 7.6 28 9z" />
  </svg>
);

export const CopyIcon = ({size = 14, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M28 10v18H10V10h18m0-2H10a2 2 0 0 0-2 2v18a2 2 0 0 0 2 2h18a2 2 0 0 0 2-2V10a2 2 0 0 0-2-2z" />
  </svg>
);

export const ArrowRight = ({size = 12, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M18 6L16.6 7.4 24.2 15H4v2h20.2l-7.6 7.6L18 26l10-10z" />
  </svg>
);

export const ArrowLeft = ({size = 12, className}: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 32 32" fill="currentColor" className={className} aria-hidden="true">
    <path d="M14 26l1.4-1.4L7.8 17H28v-2H7.8l7.6-7.6L14 6 4 16z" />
  </svg>
);
