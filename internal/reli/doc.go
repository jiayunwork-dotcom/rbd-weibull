// Package reli evaluates the reliability of a model diagram: the system
// reliability at a mission time t by folding leaf Weibull reliabilities
// through series / parallel / k-of-n combinations, an approximation of the
// system hazard rate, Birnbaum importance per leaf, and the system mean
// time to failure obtained by numerical integration of R(t).
//
// All combination rules assume independent components and evaluate every
// node at the same mission time t; intermediate nodes never substitute
// their MTTF for a reliability value. The k-of-n node uses the binomial
// sum when its children are identically distributed and an exact state
// enumeration otherwise, bounded by MaxEnumerateSubblocks.
package reli
